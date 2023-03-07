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

//Cookieに保存されているuser_idを取り出す
export function getUserId() {
  if (typeof document !== 'undefined') {
    try {
      const userId = document.cookie.split('; ')[1].split('=')[1];
      console.log(typeof (userId));
      return parseInt(userId);
    } catch (err) {
      console.log(err)
    }
  }
  return;
}
