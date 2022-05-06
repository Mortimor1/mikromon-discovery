package subnet

import (
	"context"
	"encoding/json"
	"github.com/Mortimor1/mikromon-discovery/internal/webserver/handlers"
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"github.com/gorilla/mux"
	"net/http"
)

type handler struct {
	logger *logging.Logger
	r      *SubnetRepository
}

func NewSubnetHandler(logger *logging.Logger, r *SubnetRepository) handlers.Handler {
	return &handler{
		logger: logger,
		r:      r,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/subsets", h.GetSubnets).Methods("GET")
	router.HandleFunc("/subsets/{id}", h.GetSubnetById).Methods("GET")
	router.HandleFunc("/subsets/{id}", h.CreateSubnet).Methods("POST")
	router.HandleFunc("/subsets/{id}", h.UpdateSubnet).Methods("PUT")
	router.HandleFunc("/subsets/{id}", h.DeleteSubnet).Methods("DELETE")
}

func (h *handler) GetSubnets(writer http.ResponseWriter, _ *http.Request) {
	d, err := h.r.FindAll(context.Background())
	if err != nil {
		json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
	}
	json.NewEncoder(writer).Encode(d)
}

func (h *handler) GetSubnetById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	d, err := h.r.FindOne(context.Background(), vars["id"])

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(writer).Encode(d)
}

func (h *handler) CreateSubnet(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusCreated)

	var s Subnet
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&s); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer request.Body.Close()
	id, err := h.r.Create(context.Background(), &s)

	if err != nil {
		json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
	}

	json.NewEncoder(writer).Encode(map[string]string{"id": id})
}

func (h *handler) UpdateSubnet(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
	vars := mux.Vars(request)

	var s Subnet
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&s); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer request.Body.Close()

	s.Id = vars["id"]

	errUpdate := h.r.Update(context.Background(), &s)
	if errUpdate != nil {
		json.NewEncoder(writer).Encode(map[string]string{"error": errUpdate.Error()})
	}
}

func (h *handler) DeleteSubnet(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNoContent)

	vars := mux.Vars(request)

	errDelete := h.r.Delete(context.Background(), vars["id"])
	if errDelete != nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(map[string]string{"error": errDelete.Error()})
	}
}
