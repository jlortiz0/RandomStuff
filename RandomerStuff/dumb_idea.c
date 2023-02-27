#include <curl/curl.h>
#include <stdbool.h>
#include <stdio.h>
#include <unistd.h>
#include <signal.h>

#define TEST_URL        "https://detectportal.firefox.com/canonical.html"
#define SLEEP_DUR       150
#define SHORT_SLEEP_DUR 20
#define FF_VER          "105.0"
#define USERNAME        ""
#define PASSWORD        ""
#define HOST_URL        ""

#define ENSURE(x)                                                                                  \
    if ((x)) {                                                                                     \
        perror("curl");                                                                            \
        return 1;                                                                                  \
    }

// SIGUSR1 should just skip sleeping, not terminate the program
void dummyHandler(int signum) {
    (void) signum;
}

int main(void) {
    // Running as a systemd service will not linebuf stdout by default
    setlinebuf(stdout);
    signal(SIGUSR1, &dummyHandler);
    ENSURE(curl_global_init(CURL_GLOBAL_ALL))

    CURL *pinger = curl_easy_init();
    ENSURE(pinger == NULL)
    ENSURE(curl_easy_setopt(pinger, CURLOPT_URL, TEST_URL))
    ENSURE(curl_easy_setopt(pinger, CURLOPT_NOBODY, 1))
    ENSURE(curl_easy_setopt(pinger, CURLOPT_USE_SSL, CURLUSESSL_CONTROL))

    struct curl_slist *headers = curl_slist_append(NULL, "Referer: " HOST_URL);
    headers = curl_slist_append(headers, "Origin: " HOST_URL);
    headers = curl_slist_append(headers, "Sec-Fetch-Dest: empty");
    headers = curl_slist_append(headers, "Sec-Fetch-Mode: cors");
    headers = curl_slist_append(headers, "Sec-Fetch-Site: same-origin");
    headers = curl_slist_append(headers, "Connection: keep-alive");
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "Accept-Language: en-US,en;q=0.5");
    headers = curl_slist_append(headers, "Accept: application/json, text/javascript, */*; q=0.01");
    headers = curl_slist_append(headers, "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux "
                                         "x86_64; rv:" FF_VER ") Gecko/20100101 Firefox/" FF_VER);

    while (true) {
        if (curl_easy_perform(pinger)) {
#ifdef DEBUG_VERBOSE
            puts("Attempting login...");
#endif
            char errbuf[CURL_ERROR_SIZE];
            CURL *auther = curl_easy_init();
            ENSURE(auther == NULL)
            ENSURE(curl_easy_setopt(auther, CURLOPT_ERRORBUFFER, errbuf))
            ENSURE(curl_easy_setopt(auther, CURLOPT_URL,
                        HOST_URL "/api/authRest"))
            ENSURE(curl_easy_setopt(auther, CURLOPT_POST, 1))
            ENSURE(curl_easy_setopt(auther, CURLOPT_POSTFIELDS,
                "{\"platform\":\"Linux x86_64\",\"appversion\":\"5.0 "
                "(X11)\",\"username\":\"" USERNAME "\",\"password\":\"" PASSWORD "\"}"))
            ENSURE(curl_easy_setopt(auther, CURLOPT_HTTPHEADER, headers))

            CURLcode err = curl_easy_perform(auther);
            curl_easy_cleanup(auther);
            if (err) {
                if (errbuf[0]) {
                    puts(errbuf);
                } else {
                    puts(curl_easy_strerror(err));
                }
                // This may have been due to a temporary network error
                // Try again soon, but give it some time
                sleep(SHORT_SLEEP_DUR);
                continue;
            } else {
                putchar('\n');
            }
        } else {
#ifdef DEBUG_VERBOSE
            puts("No problem!");
#endif
        }
        sleep(SLEEP_DUR);
    }
    // Unreachable leftover code
    curl_slist_free_all(headers);
    curl_easy_cleanup(pinger);
    curl_global_cleanup();
    return 0;
}
