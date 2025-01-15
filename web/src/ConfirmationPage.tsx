import { API_URL } from "./App"

export const ConfirmationPage = () => {

  const token = '';

  const handleConfirm = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: "PUT"
    })

    if (response.ok) {
      // redirect to main page
    } else {
      alert("Failed to confirm token")
    }
  }


  return (
    <div>
      <h1>Confirmation</h1>
      <button onClick={handleConfirm}>Click to confirm</button>
    </div>
  )
}