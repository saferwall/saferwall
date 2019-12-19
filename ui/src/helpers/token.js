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
    console.log(expirationDate)
    return expirationDate < new Date()
  }
export default isTokenExpired;