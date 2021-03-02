## Overview

This service provides API Create an URL-shortener service to shorten URLs.

## Application Functionalities

API Client can send a url and be returned a shortened URL. :

        ● API Client can specify an expiration time for URLs, expired URLs must return HTTP 410
        ● Input URL should be validated and respond with error if not a valid URL
        ● Regex based blacklist for URLs, urls that match the blacklist respond with an error
        ● Visiting the Shortened URLs must redirect to the original URL with a HTTP 302 redirect,404 if not found.
        ● Hit counter for shortened URLs (increment with every hit)
        ● Admin api (requiring token) to list
            ○ Short Code
            ○ Full Url
            ○ Expiry (if any)
            ○ Number of hits
        ● Above list can filter by Short Code and keyword on origin url.
    ● Admin api to delete a URL (after deletion shortened URLs must return HTTP 410)
    ● BONUS: Add a caching layer to avoid repeated database calls on popular URLs
