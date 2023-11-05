# Node interpreter

# Basic console log


```
console.log("Hello World!")

```

# Set some variables


```
const BASEURL = 'https://httpbin.org/';

```

# Async requests


```
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
