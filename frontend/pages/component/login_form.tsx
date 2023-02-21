import React from 'react'

// Propsインタフェース
interface PropsInterface {
}

// Stateインタフェース
interface StateInterface {
  name: string;
  password: string;
}

class LoginForm extends React.Component<PropsInterface, StateInterface> {
  name = "defaultUser"
  password = "defaultPassword"

  constructor(props: PropsInterface) {
    super(props)
    this.state = {
      name: '',
      password: ''
    }
    this.onClick_Submit = this.onClick_Submit.bind(this);
  }

  private onChange_Name(event:any) {
    this.setState({name: event.target.value});
  }
 
  // フォーム変更：本文
  private onChange_Password(event:any) {
    this.setState({ password: event.target.value });
    
  }

  
  // クリック：登録
  private onClick_Submit() {
    alert(this.state.name + " " + this.state.password);
  }

  

  render() {
    return (
      <div>
        <label>
          名前
          <input type='text' value={this.state.name} onChange={(e) => this.onChange_Name(e)} />
        </label><br/>
        <label>
          パスワード
          <input type="password" value={this.state.password} onChange={(e) => this.onChange_Password(e)} />
        </label>
        <button onClick={this.onClick_Submit}>ログイン</button>
        
      </div>
    )
  }
}

export default LoginForm
