# Node interpreter

# Basic console log


``` js
console.log("Hello World!")

```

# Set some variables


``` js
const BASEURL = 'https://httpbin.org/';

```

# Async requests


``` js
await fetch(`${BASEURL}post`, {
  method: 'POST',
})
  .then(response => {
    if (response.ok) {
      return response.json();
    } else {
      throw new Error(`Request failed with status code: ${response.status}`);
    }
  })
  .then(data => {
    console.log(data);
  })
  .catch(error => {
     throw new Error(`problem with fetch: ${error}`);
  });

```

# Expected Errors

Errors are supported and can be ignored, which will desplay them but not interrupt the runbook.


``` js ignore_error
throw new Error('Something went wrong but it was expected, part of the docs!');

```

# Unexpected Errors

But by default, errors will stop the runbook and cause to fail


``` js
throw new Error('Something went wrong and it was completely unexpected!');

```
