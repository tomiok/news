if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register('/static/js/serviceworker.js')
        .then(
            () => console.log("service worker registered")
        ).catch(
        err => console.log(err)
    )
}

// install button
let deferredPrompt = null;

window.addEventListener('beforeinstallprompt', function(e) {
    // Prevent Chrome 67 and earlier from automatically showing the prompt
    e.preventDefault();
    // Stash the event so it can be triggered later.
    deferredPrompt = e;
    console.log(deferredPrompt)
});

function install() {
    console.log("trying to install")
    deferredPrompt.prompt();
    // Wait for the user to respond to the prompt
    deferredPrompt.userChoice
        .then((choiceResult) => {
            if (choiceResult.outcome === 'accepted') {
                console.log('User accepted the A2HS prompt');
            } else {
                console.log('User dismissed the A2HS prompt');
            }
            deferredPrompt = null;
        });
}