'use strict';

console.log('Feed: making connection...');
const feed = new EventSource('/feed');
feed.onopen = function () {
  console.log('Feed: connection established');
};
feed.onerror = function () {
  console.log(`Feed: error. State: ${this.readyState}`);
  if (this.readyState === EventSource.CONNECTING) {
    console.log('Feed: reconnecting...');
  } else {
    console.log('Feed: fatal error...');
  }
};
feed.onmessage = function (event) {
  console.log('Feed: event received:', event.data);
};
