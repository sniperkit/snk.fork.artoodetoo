import { Component, OnInit } from '@angular/core';
import { MdDialog, MdDialogRef } from '@angular/material';

import { OptionDialog } from '../dialogs';
import { Task, Unit } from '../model';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styles: []
})
export class DashboardComponent implements OnInit {
  tasks: Task[]
  events: Unit[]

  constructor(private api: ApiService, public dialog: MdDialog) {
    api.tasks.subscribe((tasks) => this.tasks = tasks)
    api.events.subscribe((events) => this.events = events)
  }

  ngOnInit() { }

  runTask(task: Task) {
    console.log("Running", task.name)
    this.api.runTask(task).subscribe()

    console.log("Remove after refactor")
    task.running = true
  }

  stopTask(task: Task){
    console.log("Stopping", task.name)
    this.api.stopTask(task)

    console.log("Remove after refactor")
    task.running = false
  }

  deleteTask(task: Task) {
    let dialogRef = this.dialog.open(OptionDialog, {
      width: '550px',
      data: {
        title: "Delete Task?",
        message: "Would you really like to delete " + task.name
      }
    });
    dialogRef.afterClosed().subscribe(ok => {
      if (ok) {
        this.api.deleteTask(task).subscribe()
      }
    });
  }

}
