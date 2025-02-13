package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("simple_pipeline", func(t *testing.T) {
		data := loadFixture(t, "simple-pipeline.yml")
		root := expandYAML(data)
		jobs, stages, err := extractJobsAndStages(root)

		require.NoError(t, err)
		assert.Equal(t, []string{"build", "test"}, stages)
		assert.Len(t, jobs, 2)

		buildJob := findJob(jobs, "build-job")
		assert.Equal(t, "build", buildJob.Stage)
		assert.Len(t, buildJob.Rules, 1)
		assert.Equal(t, "echo \"Building\"", buildJob.Scripts[0])
	})

	t.Run("complex_rules", func(t *testing.T) {
		data := loadFixture(t, "complex-pipeline.yml")
		root := expandYAML(data)
		jobs, stages, err := extractJobsAndStages(root)

		require.NoError(t, err)
		assert.Equal(t, []string{"deploy"}, stages)
		assert.Len(t, jobs, 2)

		deployProd := findJob(jobs, "deploy-prod")
		assert.Len(t, deployProd.Rules, 2)
		assert.Equal(t, "./deploy.sh prod", deployProd.Scripts[0])
	})
}

func loadFixture(t *testing.T, name string) []byte {
	path := filepath.Join("fixtures", name)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func findJob(jobs []Job, name string) *Job {
	for _, job := range jobs {
		if job.Name == name {
			return &job
		}
	}
	return nil
}
