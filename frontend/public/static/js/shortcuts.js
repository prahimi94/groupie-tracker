document.addEventListener("keydown", function(event) {
    if (event.ctrlKey && event.key === "h") {
        event.preventDefault();
        window.location.href = "/";
    } else if (event.ctrlKey && event.key === "a") {
      event.preventDefault();
      window.location.href = "/artists";
    } else if (event.ctrlKey && event.key === "l") {
      event.preventDefault();
      window.location.href = "/locations";
    } else if (event.ctrlKey && event.key === "d") {
        event.preventDefault();
        window.location.href = "/dates";
    } else if (event.ctrlKey && event.key === "t") {
        event.preventDefault();
        window.location.href = "/tours";
    }
  });