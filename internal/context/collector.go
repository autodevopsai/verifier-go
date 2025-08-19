package context

import (
	"fmt"
	"io"
	"strings"

	"github.com/autodevops/verifier-go/internal/agent"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CollectGitContext gathers information from the local git repository using go-git.
func CollectGitContext() (agent.AgentContext, error) {
	var ctx agent.AgentContext

	repo, err := git.PlainOpen(".")
	if err != nil {
		return ctx, fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get current branch
	head, err := repo.Head()
	if err == nil {
		ctx.Branch = head.Name().Short()
	}

	// Get staged files and diff
	idx, err := repo.Storer.Index()
	if err != nil {
		return ctx, fmt.Errorf("failed to get git index: %w", err)
	}

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return ctx, fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	headTree, err := headCommit.Tree()
	if err != nil {
		return ctx, fmt.Errorf("failed to get HEAD tree: %w", err)
	}

	var stagedFiles []string
	var diffBuilder strings.Builder

	for _, entry := range idx.Entries {
		// Check if the file is modified in the index compared to HEAD
		obj, err := repo.BlobObject(entry.Hash)
		if err != nil {
			continue // Skip if blob cannot be retrieved
		}

		headEntry, err := headTree.FindEntry(entry.Name)
		if err != nil {
			// File is new in index
		
stagedFiles = append(stagedFiles, entry.Name)
		
			reader, _ := obj.Reader()
			content, _ := io.ReadAll(reader)
			diffBuilder.WriteString(fmt.Sprintf("---\n+++ b/%s\n@@ -0,0 +1,%d @@\n+%s\n", entry.Name, strings.Count(string(content), "\n")+1, string(content)))
			continue
		}

		if entry.Hash != headEntry.Hash {
			// File is modified
		
stagedFiles = append(stagedFiles, entry.Name)

			headObj, err := repo.BlobObject(headEntry.Hash)
			if err != nil {
				continue
			}

			patch, err := headObj.Patch(obj)
			if err == nil {
				diffBuilder.WriteString(patch.String())
			}
		}
	}
	ctx.Files = stagedFiles
	ctx.Diff = diffBuilder.String()

	return ctx, nil
}
