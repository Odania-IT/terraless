'use strict';

/*
  Copyright 2019 Mike Petersen Odania IT

  This lambda is deployed by Terraless
*/

// Serves all entries ending with a "/" with "/index.html"
exports.singleEntryPointHandler = async (event) => {
    const request = event.Records[0].cf.request;

    if (request.uri.match('.*/$')) {
        console.log("Request", request);
        request.uri = '/index.html';

        return request;
    }

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

    if (request.uri.match('.*/$')) {
        request.uri += 'index.html';

        return request;
    }

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

    return request;
};
