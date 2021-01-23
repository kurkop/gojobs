import { useState, Fragment, useEffect } from "react";
import { Card, Menu, Button } from "semantic-ui-react";

import { auth, authUI } from "../firebase";

import "../css/AuthForm.css";

import firebase from "firebase";


function renderLoggedIn() {
  return (
    <div className="loggedIn-wrapper">
      <h1>You are logged in!</h1>
      <div>
        <Button onClick={() => auth.signOut()} color="yellow">
          Log out
        </Button>
      </div>
    </div>
  );
}

function AuthForm() {
  const [isLogin] = useState(true);
  const [user, setUser] = useState(null);
  
  auth.onAuthStateChanged((user) => setUser(user));

  useEffect(() => {
    if (!user) {
      authUI.start(".firebaseui-auth-container", {
        signInOptions: [firebase.auth.GithubAuthProvider.PROVIDER_ID, firebase.auth.GoogleAuthProvider.PROVIDER_ID, firebase.auth.EmailAuthProvider.PROVIDER_ID],
        signInFlow: "redirect",
      });
    }
  }, [user]);

  return (
    <div className="auth-form-wrapper">
      <Card className="auth-form-card">
        <Card.Content>
          {user ? (
            renderLoggedIn()
          ) : (
            <Fragment>
              <Card.Header className="auth-form-header">GoJobs</Card.Header>
              <Menu compact secondary>
                <Card.Description className="auth-form-description">Sign in with one of the following:</Card.Description>
              </Menu>
              {isLogin ? (
                <Fragment>
                  <div className="firebaseui-auth-container"></div>
                </Fragment>
              ) : (
                <Fragment>
                  <div className="firebaseui-auth-container"></div>
                </Fragment>
              )}
            </Fragment>
          )}
        </Card.Content>
      </Card>
    </div>
  );
}

export default AuthForm;
