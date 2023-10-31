package register

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ArminGh02/go-auth-system/internal/broker"
	"github.com/ArminGh02/go-auth-system/internal/model"
	"github.com/ArminGh02/go-auth-system/internal/repository"
	"github.com/ArminGh02/go-auth-system/internal/s3"
)

type handler struct {
	users  repository.User
	broker broker.Broker
	s3     *s3.S3
}

func NewHandler(users repository.User, s3 *s3.S3, b broker.Broker) http.Handler {
	return &handler{users: users, s3: s3, broker: b}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := mux.NewRouter()

	m.HandleFunc("/register", h.register).Methods(http.MethodPost)
	m.HandleFunc("/status", h.status).Methods(http.MethodGet)

	m.ServeHTTP(w, r)
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &model.User{
		IPAddr:     r.RemoteAddr,
		Email:      r.FormValue("email"),
		Name:       r.FormValue("name"),
		NationalID: hashString(r.FormValue("national_id")),
		Status:     model.StatusPending,
	}

	file1, _, err := r.FormFile("file1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file1.Close()

	file2, _, err := r.FormFile("file2")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file2.Close()

	err = h.s3.Put(r.Context(), user.FirstImage(), file1)
	if err != nil {
		http.Error(w, "s3 put failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.s3.Put(r.Context(), user.SecondImage(), file2)
	if err != nil {
		http.Error(w, "s3 put failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.broker.Publish(r.Context(), "registrations", []byte(user.NationalID))
	if err != nil {
		http.Error(w, "broker publish failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.users.Insert(r.Context(), user)
	if err != nil {
		http.Error(w, "db insert user failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Your request for registration has been received."))
}

func (h *handler) status(w http.ResponseWriter, r *http.Request) {
	nationalID := r.URL.Query().Get("national_id")
	nationalID = hashString(nationalID)

	user, err := h.users.GetByNationalID(r.Context(), nationalID)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.IPAddr != r.RemoteAddr {
		http.Error(w, "your IP does not match the IP of this user", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{
		Status: user.Status.String(),
	})
}

func hashString(s string) string {
	hasher := sha256.New()

	hasher.Write([]byte(s))

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
