'use strict';

/*
  Copyright 2019 Mike Petersen Odania IT

  This lambda is deployed by Terraless
*/

// Serves all entries ending with a "/" with "/index.html"
exports.singleEntryPointHandler = async (event) => {
    const request = event.Records[0].cf.request;

    // If request ends with "/" serve "/index.html"
    if (request.uri.match('.*/$')) {
        request.uri = '/index.html';

        return request;
    }

    // Make sure the request ends with a "/"
    if (request.uri.match('/[^/.]+$')) {
        return {
            status: '301',
            statusDescription: 'Found',
            headers: {
                location: [{
                    key: 'Location', value: request.uri + '/',
                }],
            }
        };
    }

    // Redirect "/index.html" to "/"
    const prefixPath = request.uri.match('(.*)/index.html');
    if (prefixPath) {
        return {
            status: '301',
            statusDescription: 'Found',
            headers: {
                location: [{
                    key: 'Location', value: prefixPath[1] + '/',
                }],
            }
        };
    }

    return request;
};

// Redirect
exports.staticHandler = async (event) => {
    const request = event.Records[0].cf.request;

    // If request ends with "/" serve "/index.html"
    if (request.uri.match('.*/$')) {
        request.uri += 'index.html';

        return request;
    }

    // Redirect "/index.html" to "/"
    const prefixPath = request.uri.match('(.*)/index.html');
    if (prefixPath) {
        return {
            status: '301',
            statusDescription: 'Found',
            headers: {
                location: [{
                    key: 'Location', value: prefixPath[1] + '/',
                }],
            }
        };
    }

    // Make sure the request ends with a "/"
    if (request.uri.match('/[^/.]+$')) {
        return {
            status: '301',
            statusDescription: 'Found',
            headers: {
                location: [{
                    key: 'Location', value: request.uri + '/',
                }],
            }
        };
    }

    // Redirect to "www.DOMAIN"
    let parts = request.origin.domainName.split('.');
    if (parts === 1) {
        return {
            status: '301',
            statusDescription: 'Found',
            headers: {
                location: [{
                    key: 'Location', value: 'https://www.' + request.origin.domainName + '/' + request.uri + '/',
                }],
            }
        };
    }

    return request;
};

// Redirect to www
exports.redirectToWww = async (event) => {
    const request = event.Records[0].cf.request;

    // Redirect to "www.DOMAIN"
    let parts = request.origin.domainName.split('.');
    if (parts === 1) {
        return {
            status: '301',
            statusDescription: 'Found',
            headers: {
                location: [{
                    key: 'Location', value: 'https://www.' + request.origin.domainName + '/' + request.uri + '/',
                }],
            }
        };
    }

    return request;
};
