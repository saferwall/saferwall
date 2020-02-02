const tokenManager = {
  getTokenExpirationDate(payload){
    const jwtData = JSON.parse(atob(payload))
    if (!jwtData.exp) {
      return null
    }

    const date = new Date(0)
    date.setUTCSeconds(jwtData.exp)

    return date
  },
  isTokenExpired(payload){
    const expirationDate = this.getTokenExpirationDate(payload)
    return expirationDate < new Date()
  }
}
export default tokenManager;
