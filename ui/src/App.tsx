import { useEffect } from "react";

function App() {
  useEffect(() => {
    console.log("sending health check");
    fetch("/api/health").then((value) => {
      console.log(`response received: ${value.status}`);
    });
  });
  return (
    <>
      <h1>AI</h1>
    </>
  );
}

export default App;
