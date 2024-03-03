## Installation

To use the api module, you need to have Go installed on your system. If you don't have Go installed, you can
download it from the [official Go website](https://golang.org/dl/).

Once you have Go installed, you can get the package by running the following command:

```bash
go get github.com/ellypaws/inkbunny/api
```

## API Methods

The application provides several API methods to interact with the platform:
(note: We always need a `Credentials` object to call these methods)

| Method                                                                                                  | Description                                               |
|---------------------------------------------------------------------------------------------------------|-----------------------------------------------------------| 
| `(user *Credentials) Login() (*Credentials, error)`                                                     | Logs in a user.                                           |
| `(user *Credentials) Logout() error`                                                                    | Logs out a user.                                          |
| `(user Credentials) LoggedIn() bool`                                                                    | Checks if a user is logged in.                            |
| `(user Credentials) Request(method string, url string, body io.Reader) (*http.Request, error)`          | Sends a request to the specified URL.                     |
| `(user Credentials) Get(url *url.URL) (*http.Response, error)`                                          | Sends a GET request to the specified URL.                 |
| `(user Credentials) Post(url *url.URL, contentType string, body io.Reader) (*http.Response, error)`     | Sends a POST request to the specified URL.                |
| `(user Credentials) PostForm(url *url.URL, values url.Values) (*http.Response, error)`                  | Sends a POST request with form data to the specified URL. |
| `(user Credentials) GetWatching() ([]UsernameID, error)`                                                | Retrieves the user's watchlist. (users you're watching)   |
| `(user Credentials) GetWatchedBy(username string) ([]WatchInfo, error)`                                 | Retrieves the user's watchers. (users watching you)       |
| `(user Credentials) SubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error)` | Retrieves the details of a submission.                    |
| `(user Credentials) SubmissionFavorites(req SubmissionRequest) (SubmissionFavoritesResponse, error)`    | Retrieves the favorites of a submission.                  |
| `(user Credentials) SearchSubmissions(req SearchRequest) (SearchResponse, error)`                       | Searches for submissions.                                 |
| `(user Credentials) OwnSubmissions() (SearchResponse, error)`                                           | Retrieves the user's own submissions.                     |
| `(user Credentials) UserSubmissions(username string) (SearchResponse, error)`                           | Retrieves the submissions of a specified user.            |

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.