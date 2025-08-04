import { Card, CardBody } from "@heroui/react";
import { createFileRoute } from "@tanstack/react-router";

const component = () => {
  return (
    <Card>
      <CardBody>Login</CardBody>
    </Card>
  );
};

export const Route = createFileRoute("/login")({
  component,
});
