import { Card, CardBody } from "@heroui/react";
import { useEffect } from "react";

function App() {
  useEffect(() => {
    console.log("sending health check");
    fetch("/api/health").then((value) => {
      console.log(`response received: ${value.status}`);
    });
  });
  return (
    <Card>
      <CardBody>
        <p>Make beautiful websites regardless of your design experience.</p>
      </CardBody>
    </Card>
  );
}

export default App;
