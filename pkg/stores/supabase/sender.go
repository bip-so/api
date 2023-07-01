package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type SupapbaseRequestBody struct {
	UserID       uint64    `json:"user_id"`
	StudioUpdate bool      `json:"studio_update"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func UpdateUserSupabase(userID uint64, studioUpdated bool) {
	requestBody, err := json.Marshal(
		SupapbaseRequestBody{UserID: userID, StudioUpdate: studioUpdated, UpdatedAt: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	url := configs.GetSupabaseConfig().SupabaseBaseurl + "/rest/v1/" + "users" // users -> table name here
	req, err := http.NewRequest(http.MethodPatch, url+"?user_id=eq."+utils.String(userID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", configs.GetSupabaseConfig().SupabaseToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Println("Supabase Status code on Patch call:", resp.StatusCode)

	if resp.StatusCode == 404 {
		requestBody, err = json.Marshal(
			SupapbaseRequestBody{UserID: userID, StudioUpdate: studioUpdated, UpdatedAt: time.Now()})
		if err != nil {
			fmt.Println(err)
		}
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apiKey", configs.GetSupabaseConfig().SupabaseToken)

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Supabase Response:", string(body), "Supabase Status code:", resp.StatusCode)
}

// is_email_verified
type SupabaseRequestBodyEmailVerified struct {
	UserID          uint64    `json:"user_id"`
	IsEmailVerified bool      `json:"is_email_verified"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func UpdateUserEmailIsVerifiedSupabase(userID uint64) {
	requestBody, err := json.Marshal(
		SupabaseRequestBodyEmailVerified{UserID: userID, IsEmailVerified: true, UpdatedAt: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	url := configs.GetSupabaseConfig().SupabaseBaseurl + "/rest/v1/" + "users" // users -> table name here
	req, err := http.NewRequest(http.MethodPatch, url+"?user_id=eq."+utils.String(userID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", configs.GetSupabaseConfig().SupabaseToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Println("Supabase Status code on Patch call:", resp.StatusCode)

	if resp.StatusCode == 404 {
		requestBody, err = json.Marshal(
			SupabaseRequestBodyEmailVerified{UserID: userID, IsEmailVerified: true, UpdatedAt: time.Now()})
		if err != nil {
			fmt.Println(err)
		}
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apiKey", configs.GetSupabaseConfig().SupabaseToken)

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Supabase Response:", string(body), "Supabase Status code:", resp.StatusCode)
}

type SupabaseRequestBodyDefaultStudio struct {
	UserID        uint64    `json:"user_id"`
	DefaultStudio uint64    `json:"default_studio"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func UpdateUserDefaultStudioSupabase(userID uint64, studioID uint64) {
	requestBody, err := json.Marshal(
		SupabaseRequestBodyDefaultStudio{UserID: userID, DefaultStudio: studioID, UpdatedAt: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	url := configs.GetSupabaseConfig().SupabaseBaseurl + "/rest/v1/" + "users" // users -> table name here
	req, err := http.NewRequest(http.MethodPatch, url+"?user_id=eq."+utils.String(userID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", configs.GetSupabaseConfig().SupabaseToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Println("Supabase Status code on Patch call:", resp.StatusCode)
	if resp.StatusCode == 404 {
		requestBody, err = json.Marshal(
			SupabaseRequestBodyDefaultStudio{UserID: userID, DefaultStudio: studioID, UpdatedAt: time.Now()})
		if err != nil {
			fmt.Println(err)
		}
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apiKey", configs.GetSupabaseConfig().SupabaseToken)

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Supabase Response:", string(body), "Supabase Status code:", resp.StatusCode)
}
