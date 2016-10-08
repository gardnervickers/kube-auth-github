package main

type TokenReviewSpec struct {
	Token string `json:"token"`
}

type TokenReviewStatus struct {
	Authenticated bool                   `json:"authenticated"`
	User          *TokenReviewStatusUser `json:"user"`
}

type TokenReviewStatusUser struct {
	Username string   `json:"user"`
	Uid      string   `json:"uid"`
	Groups   []string `json:"groups"`
}

type TokenReview struct {
	ApiVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Spec       *TokenReviewSpec   `json:"spec,omitempty"`
	Status     *TokenReviewStatus `json:"status"`
}
