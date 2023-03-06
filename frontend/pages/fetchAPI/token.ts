//Cookieに保存されているjwtTokenを取り出す
export function getToken() {
  const jwtToken = "";
  if (typeof document !== 'undefined') {
    try {
      const jwtToken = document.cookie.split('; ')[0].split('=')[1];
      return jwtToken;
    } catch (err) {
      console.log(err)
    }
  }
  return jwtToken;
}
