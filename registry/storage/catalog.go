package storage

import (
	"errors"
	"io"
	"sort"

	"github.com/docker/distribution/context"
	//"github.com/docker/distribution/registry/storage/driver"
)

// ErrFinishedWalk is used when the called walk function no longer wants
// to accept any more values.  This is used for pagination when the
// required number of repos have been found.
var ErrFinishedWalk = errors.New("finished walk")

// Returns a list, or partial list, of repositories in the registry.
// Because it's a quite expensive operation, it should only be used when building up
// an initial set of repositories.
func (reg *registry) Repositories(ctx context.Context, repos []string, last string) (n int, errVal error) {
	var foundRepos []string

	if len(repos) == 0 {
		return 0, errors.New("no space in slice")
	}

	root, err := pathFor(repositoriesRootPathSpec{})
	if err != nil {
		return 0, err
	}

	/*
	err = Walk(ctx, reg.blobStore.driver, root, func(fileInfo driver.FileInfo) error {
		filePath := fileInfo.Path()

		// lop the base path off
		repoPath := filePath[len(root)+1:]

		_, file := path.Split(repoPath)
		if file == "_layers" {
			repoPath = strings.TrimSuffix(repoPath, "/_layers")
			if repoPath > last {
				foundRepos = append(foundRepos, repoPath)
			}
			return ErrSkipDir
		} else if strings.HasPrefix(file, "_") {
			return ErrSkipDir
		}

		// if we've filled our array, no need to walk any further
		if len(foundRepos) == len(repos) {
			return ErrFinishedWalk
		}

		return nil
	})
	*/

	catalogRoot := "/catalog" + root
	children, err := reg.blobStore.driver.List(ctx, catalogRoot)
	if err != nil {
		return 0, err
	}

	for _, child := range children {
		// lop the base path off
		repoPath := child[len(root)+1:]

		if repoPath > last {
			foundRepos = append(foundRepos, repoPath)
		}
	}
	sort.Strings(foundRepos)

	n = copy(repos, foundRepos)

	// Signal that we have no more entries by setting EOF
	if len(foundRepos) <= len(repos) && err != ErrFinishedWalk {
		errVal = io.EOF
	}

	return n, errVal
}
