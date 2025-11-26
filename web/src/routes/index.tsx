import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div style={{ padding: '2rem' }}>
      <h1>Nexus</h1>
    </div>
  )
}
