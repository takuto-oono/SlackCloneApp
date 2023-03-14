//Cookieに保存されているjwtTokenを取り出す
export function getToken(): string{
  if (typeof document !== 'undefined') {
    try {
      const cookie = document.cookie.split('; ').find(row => row.startsWith('token'));
      if (cookie !== undefined) {
        const jwtToken = cookie.split('=')[1];
        return jwtToken;
      }
    } catch (err) {
      console.log(err)
      return "";
    }
  }
  return "";
}

//Cookieに保存されているuser_idを取り出す
export function getUserId() {
  if (typeof document !== 'undefined') {
    try {
      const cookie = document.cookie.split('; ').find(row => row.startsWith('user_id'));
      if (cookie !== undefined) {
        const user_id = parseInt(cookie.split('=')[1]);
        return user_id;
      }
    } catch (err) {
      console.log(err)
      return;
    }
  }
}
