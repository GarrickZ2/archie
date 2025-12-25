package clone

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

// CopyManager handles file and directory copying
type CopyManager struct {
	fs afero.Fs
}

// CopyOptions defines options for copying
type CopyOptions struct {
	SourcePath string
	TargetPath string
	Items      []CloneItem
	Overwrite  bool
}

// CopyResult contains the result of a copy operation
type CopyResult struct {
	CopiedFiles  []string
	CopiedDirs   []string
	SkippedFiles []string
	Errors       []error
}

// NewCopyManager creates a new copy manager
func NewCopyManager(fs afero.Fs) *CopyManager {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &CopyManager{fs: fs}
}

// Copy performs the copying operation based on options
func (c *CopyManager) Copy(opts CopyOptions) (*CopyResult, error) {
	result := &CopyResult{
		CopiedFiles:  []string{},
		CopiedDirs:   []string{},
		SkippedFiles: []string{},
		Errors:       []error{},
	}

	for _, item := range opts.Items {
		srcPath := filepath.Join(opts.SourcePath, item.Name)
		dstPath := filepath.Join(opts.TargetPath, item.Name)

		// Check if source exists
		exists, err := afero.Exists(c.fs, srcPath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to check %s: %w", item.Name, err))
			continue
		}
		if !exists {
			// Skip if source doesn't exist
			continue
		}

		if item.Type == "file" {
			if err := c.copyFile(srcPath, dstPath, opts.Overwrite); err != nil {
				if _, ok := err.(*SkippedError); ok {
					result.SkippedFiles = append(result.SkippedFiles, item.Name)
				} else {
					result.Errors = append(result.Errors, fmt.Errorf("failed to copy file %s: %w", item.Name, err))
				}
			} else {
				result.CopiedFiles = append(result.CopiedFiles, item.Name)
			}
		} else if item.Type == "directory" {
			if err := c.copyDirectory(srcPath, dstPath, opts.Overwrite, result); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to copy directory %s: %w", item.Name, err))
			} else {
				result.CopiedDirs = append(result.CopiedDirs, item.Name)
			}
		}
	}

	return result, nil
}

// SkippedError indicates a file was skipped
type SkippedError struct {
	Path string
}

func (e *SkippedError) Error() string {
	return fmt.Sprintf("skipped existing file: %s", e.Path)
}

// copyFile copies a single file from source to target
func (c *CopyManager) copyFile(srcPath, dstPath string, overwrite bool) error {
	// Check if destination exists
	exists, err := afero.Exists(c.fs, dstPath)
	if err != nil {
		return err
	}

	if exists && !overwrite {
		return &SkippedError{Path: dstPath}
	}

	// Open source file
	srcFile, err := c.fs.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := c.fs.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy content
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	// Get source file info for permissions
	srcInfo, err := c.fs.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	// Apply same permissions to destination
	if err := c.fs.Chmod(dstPath, srcInfo.Mode()); err != nil {
		// Non-fatal: just log
		return nil
	}

	return nil
}

// copyDirectory recursively copies a directory
func (c *CopyManager) copyDirectory(srcDir, dstDir string, overwrite bool, result *CopyResult) error {
	// Create target directory
	if err := c.fs.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read source directory entries
	entries, err := afero.ReadDir(c.fs, srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Recursively copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			// Recursive call for subdirectories
			if err := c.copyDirectory(srcPath, dstPath, overwrite, result); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := c.copyFile(srcPath, dstPath, overwrite); err != nil {
				if _, ok := err.(*SkippedError); ok {
					relPath, _ := filepath.Rel(dstDir, dstPath)
					result.SkippedFiles = append(result.SkippedFiles, relPath)
				} else {
					return err
				}
			} else {
				relPath, _ := filepath.Rel(dstDir, dstPath)
				result.CopiedFiles = append(result.CopiedFiles, relPath)
			}
		}
	}

	return nil
}

// shouldCopy determines if a file should be copied based on existence and overwrite flag
func (c *CopyManager) shouldCopy(dstPath string, overwrite bool) (bool, error) {
	exists, err := afero.Exists(c.fs, dstPath)
	if err != nil {
		return false, err
	}

	// Copy if doesn't exist or overwrite is true
	return !exists || overwrite, nil
}
