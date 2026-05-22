package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type Contact struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Adress      string `json:"adress"`
	PhoneNumber string `json:"phoneNumber"`
}

type contactStore struct {
	mu       sync.RWMutex
	contacts []Contact
	nextID   int
}

func newContactStore() *contactStore {
	return &contactStore{contacts: make([]Contact, 0), nextID: 1}
}

func (s *contactStore) list() []Contact {
	s.mu.RLock()
	defer s.mu.RUnlock()

	contacts := make([]Contact, len(s.contacts))
	copy(contacts, s.contacts)
	return contacts
}

func (s *contactStore) get(id int) (Contact, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, contact := range s.contacts {
		if contact.ID == id {
			return contact, true
		}
	}

	return Contact{}, false
}

func (s *contactStore) create(contact Contact) Contact {
	s.mu.Lock()
	defer s.mu.Unlock()

	contact.ID = s.nextID
	s.nextID++
	s.contacts = append(s.contacts, contact)
	return contact
}

func (s *contactStore) update(id int, updated Contact) (Contact, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for index := range s.contacts {
		if s.contacts[index].ID == id {
			updated.ID = id
			s.contacts[index] = updated
			return updated, true
		}
	}

	return Contact{}, false
}

func (s *contactStore) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for index := range s.contacts {
		if s.contacts[index].ID == id {
			s.contacts = append(s.contacts[:index], s.contacts[index+1:]...)
			return true
		}
	}

	return false
}

func main() {
	store := newContactStore()
	router := mux.NewRouter()

	router.HandleFunc("/contacts", listContactsHandler(store)).Methods(http.MethodGet)
	router.HandleFunc("/contacts/{id}", getContactHandler(store)).Methods(http.MethodGet)
	router.HandleFunc("/contacts", createContactHandler(store)).Methods(http.MethodPost)
	router.HandleFunc("/contacts/{id}", updateContactHandler(store)).Methods(http.MethodPut)
	router.HandleFunc("/contacts/{id}", deleteContactHandler(store)).Methods(http.MethodDelete)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("server error:", err)
	}
}

func listContactsHandler(store *contactStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		respondJSON(writer, http.StatusOK, store.list())
	}
}

func getContactHandler(store *contactStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := parseID(request)
		if err != nil {
			respondError(writer, http.StatusBadRequest, err.Error())
			return
		}

		contact, found := store.get(id)
		if !found {
			respondError(writer, http.StatusNotFound, "contact not found")
			return
		}

		respondJSON(writer, http.StatusOK, contact)
	}
}

func createContactHandler(store *contactStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		contact, err := decodeContact(request)
		if err != nil {
			respondError(writer, http.StatusBadRequest, err.Error())
			return
		}

		created := store.create(contact)
		respondJSON(writer, http.StatusCreated, created)
	}
}

func updateContactHandler(store *contactStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := parseID(request)
		if err != nil {
			respondError(writer, http.StatusBadRequest, err.Error())
			return
		}

		contact, err := decodeContact(request)
		if err != nil {
			respondError(writer, http.StatusBadRequest, err.Error())
			return
		}

		updated, found := store.update(id, contact)
		if !found {
			respondError(writer, http.StatusNotFound, "contact not found")
			return
		}

		respondJSON(writer, http.StatusOK, updated)
	}
}

func deleteContactHandler(store *contactStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := parseID(request)
		if err != nil {
			respondError(writer, http.StatusBadRequest, err.Error())
			return
		}

		if !store.delete(id) {
			respondError(writer, http.StatusNotFound, "contact not found")
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

func decodeContact(request *http.Request) (Contact, error) {
	var payload Contact
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		return Contact{}, errors.New("invalid request body")
	}

	if strings.TrimSpace(payload.Name) == "" {
		return Contact{}, errors.New("name is required")
	}
	if strings.TrimSpace(payload.Adress) == "" {
		return Contact{}, errors.New("adress is required")
	}
	if strings.TrimSpace(payload.PhoneNumber) == "" {
		return Contact{}, errors.New("phoneNumber is required")
	}

	return Contact{
		Name:        strings.TrimSpace(payload.Name),
		Adress:      strings.TrimSpace(payload.Adress),
		PhoneNumber: strings.TrimSpace(payload.PhoneNumber),
	}, nil
}

func parseID(request *http.Request) (int, error) {
	params := mux.Vars(request)
	rawID := params["id"]
	if strings.TrimSpace(rawID) == "" {
		return 0, errors.New("id is required")
	}

	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		return 0, errors.New("id must be a positive integer")
	}

	return id, nil
}

func respondJSON(writer http.ResponseWriter, statusCode int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(payload)
}

func respondError(writer http.ResponseWriter, statusCode int, message string) {
	respondJSON(writer, statusCode, map[string]string{"error": message})
}
