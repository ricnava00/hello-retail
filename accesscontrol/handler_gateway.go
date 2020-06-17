package accesscontrol

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Wrap(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config, err := getConfig()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO: The token in the authorization header
		//       should be swapped for a policy.
		policy := r.Header.Get("Authorization")

		function, ok := mux.Vars(r)["name"]
		if !ok {
			fmt.Println("could not get function name")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if policy == "service" {
			err = checkPermissionsWithHeader(config, function, r)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err = checkPermissionsWithGraphs(config, function, policy, r)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = wrapRequestBody(config, policy, r)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Header.Set("Authorization", "service")
		}

		handler.ServeHTTP(w, r)
	}
}
