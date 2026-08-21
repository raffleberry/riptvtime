package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/glebarez/sqlite"
	"github.com/raffleberry/riptvtime/internal/api"
	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/utils"
	glog "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	Data string
}

var _db *gorm.DB
var BrowserOpened = false

func conn() *gorm.DB {
	if _db == nil {
		osCfgDir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		cfgDir := filepath.Join(osCfgDir, "riptvtime")

		err = os.MkdirAll(cfgDir, 0755)
		if err != nil {
			panic(err)
		}

		dbPath := filepath.Join(cfgDir, "config.db")
		slog.Debug("Initializing Config Sqlite Database", "path", dbPath)

		_db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%v?", dbPath)), &gorm.Config{
			Logger: glog.NewSlogLogger(slog.Default(), glog.Config{
				IgnoreRecordNotFoundError: true,
			}),
		})
		if err != nil {
			panic(err)
		}

		err = _db.AutoMigrate(&Config{})

		if err != nil {
			slog.Error("Failed to migrate", "err", err)
			panic("Failed to migrate database")
		}
	}
	return _db
}

func getConfigIfExists() *config.Config {
	cRow := Config{}

	err := conn().Last(&cRow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		panic(err)
	}
	c := config.Config{}
	err = json.Unmarshal([]byte(cRow.Data), &c)
	if err != nil {
		// bad data, but lets get a newer one from the user.
		slog.Warn("Bad data in config", "err", err, "data", cRow.Data)
		return nil
	}

	return &c
}

func GetConfigFromUser() (*config.Config, error) {
	rv := getConfigIfExists()
	if rv != nil {
		return config.LoadFromUISetup(rv)
	}

	done := make(chan bool)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(landingPage))
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		tmdbApiKey := r.FormValue("TmdbApiKey")
		ip := r.FormValue("Ip")
		portStr := r.FormValue("Port")
		maxRetriesStr := r.FormValue("TmdbMaxRetries")

		port, err := strconv.Atoi(portStr)
		if err != nil {
			http.Error(w, "invalid port", http.StatusBadRequest)
			return
		}
		maxRetries, err := strconv.Atoi(maxRetriesStr)
		if err != nil {
			http.Error(w, "invalid maxRetries", http.StatusBadRequest)
			return
		}

		rv = &config.Config{
			TmdbApiKey:     tmdbApiKey,
			Ip:             ip,
			Port:           port,
			TmdbMaxRetries: maxRetries,
		}

		json, err := json.Marshal(rv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = conn().Create(&Config{
			Data: string(json),
		}).Error

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, "http://%s:%d/", rv.Ip, rv.Port)
		if err != nil {
			panic(err)
		}

		done <- true
	})

	url := "http://127.0.0.1:5667"
	s := api.NewServer("127.0.0.1:5667", mux)
	err := s.Start()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = s.Stop()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("=====================")
	fmt.Printf("App setup page: %v\n", url)
	fmt.Println("=====================")

	err = utils.OpenBrowser(url)
	BrowserOpened = true
	if err != nil {
		panic(err)
	}

	<-done
	slog.Info("Done")

	return config.LoadFromUISetup(rv)
}

const landingPage string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome</title>
    <style>
        * {
            box-sizing: border-box;
        }

        body {
            margin: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
            background: #f5f7fb;
            color: #1f2937;
        }

        .card {
            width: min(440px, calc(100% - 32px));
            padding: 32px;
            background: white;
            border-radius: 14px;
            box-shadow: 0 10px 35px rgba(0, 0, 0, 0.08);
        }

        h1 {
            margin: 0 0 8px;
            font-size: 28px;
        }

        .subtitle {
            margin: 0 0 28px;
            color: #6b7280;
            line-height: 1.5;
        }

        .field {
            margin-bottom: 18px;
        }

        label {
            display: block;
            margin-bottom: 7px;
            font-size: 14px;
            font-weight: 600;
        }

        input {
            width: 100%;
            padding: 11px 12px;
            border: 1px solid #d1d5db;
            border-radius: 8px;
            font-size: 15px;
            outline: none;
            transition: border-color 0.15s, box-shadow 0.15s;
        }

        input:focus {
            border-color: #6366f1;
            box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.12);
        }

        button {
            width: 100%;
            padding: 12px;
            margin-top: 6px;
            border: 0;
            border-radius: 8px;
            background: #4f46e5;
            color: white;
            font-size: 15px;
            font-weight: 600;
            cursor: pointer;
        }

        button:hover {
            background: #4338ca;
        }

        button:disabled {
            opacity: 0.6;
            cursor: wait;
        }

        #error {
            display: none;
            margin-bottom: 18px;
            padding: 11px 13px;
            border-radius: 8px;
            background: #fef2f2;
            color: #b91c1c;
            border: 1px solid #fecaca;
            font-size: 14px;
        }

        #success {
            display: none;
            margin-bottom: 18px;
            padding: 11px 13px;
            border-radius: 8px;
            background: #f0fdf4;
            color: #15803d;
            border: 1px solid #bbf7d0;
            font-size: 14px;
        }
    </style>
</head>

<body>
    <main class="card">
        <h1>Welcome 👋</h1>
        <p class="subtitle">
            Let's get things configured. Enter the details below to get started.
        </p>

        <div id="error"></div>
        <div id="success"></div>

        <form id="config-form">
            <div class="field">
                <label for="tmdbApiKey">TMDB API Key</label>
                <input
                    id="tmdbApiKey"
                    name="TmdbApiKey"
                    type="password"
                    placeholder="Enter your TMDB API key"
                    required
                >
				<div class="hint">
    				To register for an API key, click the <a href="https://www.themoviedb.org/settings/api" target="_blank" rel="noopener noreferrer">API link</a> from within your account settings page.
				</div>
            </div>

            <div class="field">
                <label for="ip">IP Address</label>
                <input
                    id="ip"
                    name="Ip"
                    type="text"
                    value="127.0.0.1"
                    required
                >
            </div>

            <div class="field">
                <label for="port">Port</label>
                <input
                    id="port"
                    name="Port"
                    type="number"
                    value="5667"
                    min="1"
                    max="65535"
                    required
                >
            </div>

            <div class="field">
                <label for="tmdbMaxRetries">TMDB Max Retries</label>
                <input
                    id="tmdbMaxRetries"
                    name="TmdbMaxRetries"
                    type="number"
                    value="10"
                    min="0"
                    required
                >
            </div>

            <button id="submit" type="submit">
                Continue
            </button>
        </form>
    </main>

    <script>
        const form = document.getElementById('config-form');
        const error = document.getElementById('error');
        const success = document.getElementById('success');
        const submit = document.getElementById('submit');

        function showError(message) {
            error.textContent = message;
            error.style.display = 'block';
            success.style.display = 'none';
        }

        const loadrtt = (afterMs, url) => {
            setTimeout(() => {
                window.location.href = url;
            }, afterMs);
        }

        form.addEventListener('submit', async function (event) {
            event.preventDefault();

            error.style.display = 'none';
            success.style.display = 'none';
            submit.disabled = true;

            try {
                const response = await fetch(window.location.href, {
                    method: 'POST',
                    body: new URLSearchParams(new FormData(form))
                });

                const address = await response.text();

                if (!response.ok) {
                    showError('Failed to save configuration.');
                    return;
                }

                success.textContent = 'Configuration saved successfully. Reloading in 3 seconds...';
                success.style.display = 'block';
                form.style.display = 'none';

                loadrtt(3000, address);
            } catch (err) {
                console.log(err)
                showError('Unable to connect to the server.');
            } finally {
                submit.disabled = false;
            }
        });
    </script>
</body>
</html>
`
