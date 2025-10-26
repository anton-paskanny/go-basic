package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"order-api-cart/models"
)

// AuthServiceClient handles communication with the auth service
type AuthServiceClient struct {
	baseURL string
	client  *http.Client
}

// NewAuthServiceClient creates a new auth service client
func NewAuthServiceClient(baseURL string) *AuthServiceClient {
	return &AuthServiceClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUserByID fetches user data from the auth service
func (c *AuthServiceClient) GetUserByID(userID string) (*models.ExternalUser, error) {
	url := fmt.Sprintf("%s/users/%s", c.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user not found: status %d", resp.StatusCode)
	}

	var user models.ExternalUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// ValidateUser validates that a user exists in the auth service
func (c *AuthServiceClient) ValidateUser(userID string) error {
	_, err := c.GetUserByID(userID)
	return err
}

// ProductServiceClient handles communication with the product service
type ProductServiceClient struct {
	baseURL string
	client  *http.Client
}

// NewProductServiceClient creates a new product service client
func NewProductServiceClient(baseURL string) *ProductServiceClient {
	return &ProductServiceClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetProductByID fetches product data from the product service
func (c *ProductServiceClient) GetProductByID(productID string) (*models.ExternalProduct, error) {
	url := fmt.Sprintf("%s/products/%s", c.baseURL, productID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found: status %d", resp.StatusCode)
	}

	var product models.ExternalProduct
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &product, nil
}

// UpdateProductQuantity updates product quantity in the product service
func (c *ProductServiceClient) UpdateProductQuantity(productID string, quantityChange int) error {
	url := fmt.Sprintf("%s/products/%s/quantity", c.baseURL, productID)

	payload := map[string]int{
		"change": quantityChange,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update product quantity: status %d", resp.StatusCode)
	}

	return nil
}
