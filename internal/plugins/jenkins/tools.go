package jenkins

import (
	"context"
	"fmt"
)

// listJobs lists all Jenkins jobs
func (p *Plugin) listJobs(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	jobs, err := p.client.GetJobs(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"count": len(jobs),
		"jobs":  jobs,
	}, nil
}

// getJob gets detailed information about a Jenkins job
func (p *Plugin) getJob(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	jobName, ok := params["job_name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'job_name' parameter")
	}

	job, err := p.client.GetJob(ctx, jobName)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// triggerBuild triggers a build for a Jenkins job
func (p *Plugin) triggerBuild(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	jobName, ok := params["job_name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'job_name' parameter")
	}

	// Extract build parameters if provided
	var buildParams map[string]string
	if paramsRaw, exists := params["parameters"]; exists {
		if paramsMap, ok := paramsRaw.(map[string]interface{}); ok {
			buildParams = make(map[string]string)
			for k, v := range paramsMap {
				buildParams[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	err := p.client.TriggerBuild(ctx, jobName, buildParams)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Build triggered for job: %s", jobName),
		"job":     jobName,
	}, nil
}

// getBuild gets detailed information about a specific build
func (p *Plugin) getBuild(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	jobName, ok := params["job_name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'job_name' parameter")
	}

	buildNumber, ok := params["build_number"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'build_number' parameter")
	}

	build, err := p.client.GetBuild(ctx, jobName, int(buildNumber))
	if err != nil {
		return nil, err
	}

	return build, nil
}

// getBuildLog gets console log for a build
func (p *Plugin) getBuildLog(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	jobName, ok := params["job_name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'job_name' parameter")
	}

	buildNumber, ok := params["build_number"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'build_number' parameter")
	}

	start := 0
	if startRaw, exists := params["start"]; exists {
		if startFloat, ok := startRaw.(float64); ok {
			start = int(startFloat)
		}
	}

	log, err := p.client.GetBuildLog(ctx, jobName, int(buildNumber), start)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"job":          jobName,
		"build_number": int(buildNumber),
		"log":          log,
		"length":       len(log),
	}, nil
}

// listBuilds lists build history for a Jenkins job
func (p *Plugin) listBuilds(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	jobName, ok := params["job_name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'job_name' parameter")
	}

	limit := 0
	if limitRaw, exists := params["limit"]; exists {
		if limitFloat, ok := limitRaw.(float64); ok {
			limit = int(limitFloat)
		}
	}

	builds, err := p.client.GetBuilds(ctx, jobName, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"job":    jobName,
		"count":  len(builds),
		"builds": builds,
	}, nil
}
