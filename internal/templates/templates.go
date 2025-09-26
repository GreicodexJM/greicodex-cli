package templates

import "grei-cli/internal/core/recipe"

type Data struct {
	recipe.Recipe
	Year int
}
