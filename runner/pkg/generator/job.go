package generator

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	collector   = "staging-collector.newrelic.com"
	driverImage = "quay.io/emiliogarcia_1/traffic-driver:latest"
	secretName  = "app-secret"
)

func (run *Run) CreateJob() batchv1.Job {
	backoffLimit := int32(1)
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-job", run.Name),
			Labels: map[string]string{
				"app": run.Name,
			},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": run.Name,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  fmt.Sprintf("%s-app", run.Name),
							Image: run.App.Image,
							Env:   run.appEnv(),
							EnvFrom: []v1.EnvFromSource{
								{
									SecretRef: &v1.SecretEnvSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: secretName,
										},
									},
								},
							},
						},
						{
							Name:  fmt.Sprintf("%s-driver", run.Name),
							Image: driverImage,
							Env:   run.driverEnv(),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}
}

func (run *Run) appEnv() []v1.EnvVar {
	vars := []v1.EnvVar{
		{
			Name:  "NEW_RELIC_HOST",
			Value: collector,
		},
		{
			Name:  "NEW_RELIC_APP_NAME",
			Value: run.Name,
		},
	}

	for k, v := range run.App.EnvVars {
		vars = append(vars, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	return vars
}

func (run *Run) driverEnv() []v1.EnvVar {
	vars := []v1.EnvVar{
		{
			Name:  "TRAFFIC_DRIVER_DELAY",
			Value: fmt.Sprint(run.TrafficDriver.Delay),
		},
		{
			Name:  "SERVICE_ENDPOINT",
			Value: run.TrafficDriver.Endpoint,
		},
		{
			Name:  "CONCURRENT_REQUESTS",
			Value: fmt.Sprint(run.TrafficDriver.Traffic.Users),
		},
		{
			Name:  "REQUESTS_PER_SECOND",
			Value: fmt.Sprint(run.TrafficDriver.Traffic.Rate),
		},
		{
			Name:  "DURATION",
			Value: fmt.Sprintf("%ds", run.TrafficDriver.Traffic.Duration),
		},
	}
	return vars
}
