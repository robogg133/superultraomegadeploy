package routes

import (
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/pkg/kube"
	"github.com/rs/zerolog/log"
)

func PodsList(k *kube.Kube) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default"
		}
		list, err := k.ListPods(r.Context(), namespace)
		if err != nil {
			log.Ctx(r.Context()).Err(err).Send()
			response.New().InternalServerError(w)
			return
		}
		resp := make([]map[string]any, 0)

		for _, pod := range list {
			r := make(map[string]any)
			r["name"] = pod.Name
			r["namespace"] = pod.Namespace
			r["creation_time"] = pod.CreationTimestamp.Time
			resp = append(resp, r)
		}
		log.Ctx(r.Context()).Debug().Int("len_pods", len(resp)).Send()

		response.New().Response(resp).Send(w, 200)
	}
}
