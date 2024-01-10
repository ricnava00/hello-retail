'use strict';

const got = require('got');
const { KV_Store } = require('kv-store');

const constants = {
  MODULE: 'authenticate.js',
  AUTHENTICATION_URL: process.env.AUTHENTICATION_URL
};

const impl = {
  getEmail: (event, callback) => {
    if("authorization" in event.headers)
    {
      got.get(constants.AUTHENTICATION_URL, {headers:{"Authorization": event.headers.authorization}}).then(r=>{
	callback(null,{user:JSON.parse(r.body)['email']});
      }).catch(r=>{
        console.log(r);
        callback(r);
      })
    }
    else
    {
      callback("Unauthorized")
    }
  },
};

module.exports = (event, context, callback) => {
  console.log(JSON.stringify(event));
  impl.getEmail(event, callback);
};