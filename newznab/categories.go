package newznab

// Category describes a newznab category
type Category int

// Newznab category constants
const (
	// TV Categories
	// CategoryTVAll is for all shows
	CategoryTVAll Category = 5000
	// CategoryTVForeign is for foreign shows
	CategoryTVForeign Category = 5020
	// CategoryTVSD is for standard-definition shows
	CategoryTVSD Category = 5030
	// CategoryTVHD is for high-definition shows
	CategoryTVHD Category = 5040
	// CategoryTVOther is for other shows
	CategoryTVOther Category = 5050
	// CategoryTVSport is for sports shows
	CategoryTVSport Category = 5060

	// Movie categories
	// CategoryMovieAll is for all movies
	CategoryMovieAll Category = 2000
	// CategoryMovieForeign is for foreign movies
	CategoryMovieForeign Category = 2010
	// CategoryMovieOther is for other movies
	CategoryMovieOther Category = 2020
	// CategoryMovieSD is for standard-definition movies
	CategoryMovieSD Category = 2030
	// CategoryMovieHD is for high-definition movies
	CategoryMovieHD Category = 2040
	// CategoryMovieBluRay is for blu-ray movies
	CategoryMovieBluRay Category = 2050
	// CategoryMovie3D is for 3-D movies
	CategoryMovie3D Category = 2060
)
