const getTokenExpirationDate = (payload) => {
    const jwtData = JSON.parse(atob(payload))
    if (!jwtData.exp) {
      return null
    }

    const date = new Date(0)
    date.setUTCSeconds(jwtData.exp)

    return date
  }
const isTokenExpired = (payload) => {
    const expirationDate = getTokenExpirationDate(payload)
    return expirationDate < new Date()
  }
export default isTokenExpired;