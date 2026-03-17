export async function getMock<T>(path: string): Promise<T> {
  const response = await fetch(`/mocks/${path}`)

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Mock HTTP ${response.status}: ${errorText}`)
  }

  return response.json() as Promise<T>
}
