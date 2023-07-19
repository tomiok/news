const CACHE_NAME = "web6am";
const urlsToCache = [
    "/static/favicon.ico",
    "/static/styles.css"
];

self.addEventListener("install", e => {
    e.waitUntil(
        caches.open(CACHE_NAME).then(cache => {
            return cache.addAll(urlsToCache);
        }).catch(
            err => console.log(err)
        )
    )
});

// Cache First Policy
self.addEventListener("fetch", event => {
    event.respondWith(
        caches.match(event.request)
            .then(function(response) {
                if (response) {  // CACHE HIT
                    return response;
                } else {    // CACHE MISS
                    return fetch(event.request);
                }
            })
    );
});