import { Card, CardBody } from "@heroui/react";
import { createFileRoute } from "@tanstack/react-router";

const component = () => {
  // TODO: check credentials, route to login or dashboard depending on result
  return (
    <Card>
      <CardBody>Hello</CardBody>
    </Card>
  );
};

export const Route = createFileRoute("/")({
  component,
});
