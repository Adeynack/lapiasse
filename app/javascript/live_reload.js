if (process.env.NODE_ENV === "development") {
  const source = new EventSource("http://localhost:5678");
  source.onmessage = (event) => {
    if (event.data === "reload") {
      window.location.reload();
    }
  };
}
