import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ApiService } from '../api.service';
import { MdSnackBar } from '@angular/material';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styles: []
})
export class LoginComponent implements OnInit {

  password: string
  username: string
  error: string

  constructor(private api: ApiService, private router: Router, private snackBar: MdSnackBar) { }

  ngOnInit() {
  }

  private loginSuccess() {
    this.api.getAll()
    this.router.navigateByUrl('/')
  }

  private loginError() {
    this.snackBar.open("Invalid credentials", "", {duration: 2000, extraClasses: ["snackbar-error"]})
    this.error = "Invalid credentials"
  }

  login() {
    console.info("Login")
    let self = this
    this.api.login(this.username, this.password).subscribe(
      () => self.loginSuccess(),
      () => self.loginError()
    )
  }

}
