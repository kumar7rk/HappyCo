const {Client} = require('pg')
var express = require('express');
var app = express();
require('dotenv').config();

// object for all the inspection data
var inspection_data = {
  i_id:[],
  i_folder_id:[],
  i_created_at: []
};

// object for all the report data
var report_data = {
 r_id:[],
 r_created_at:[],
};

// object for iap
var iap_data = {
  iap_expired_at:[]
};

// object for business notes from admin
var business_notes_data = {
  b_body:[]
};

//app.listen(3000, () => console.log('Happy app listening on port 3000'));

// forming db url
const client = new Client({

  user: process.env.DB_USER,
  host: process.env.DB_HOST,
  database: process.env.DB_NAME,
  password: process.env.DB_PASS
})

// connecting to db
// if and when the connections is established call the userdata method for now--> should be called in new 
client.connect((err) => {
  if (err) {
    console.error('connection error', err.stack)
  } else {
    console.log('connected')
    //console.log(process.env)
    getUserData(65135);// ideally call in new conversation method
  }
});

/* Method to be called when a new message is received in intercom
* calls other methods which sets the value of the objects
*/
function getUserData(u_id){
//  console.log("getUserData");
  
  getLatestInspections(u_id);
  getLatestReports(u_id);
  getIAPStatus(u_id);
  getBusinessNotes(u_id);
  
  //waiting for 2 second to print values on console
  //only used for testing should be removed in the release version
  setTimeout(printValues,2000);
}

// prints the values of objects on console
// and calls noteBuilder--> should be called in new conversation

function printValues(){
  /*console.log(inspection_data);
  console.log(report_data);
  console.log(iap_data);
  console.log(business_notes_data);*/

  noteBuilder();//should be called in new conversation

}

// gets the data (folder_id, inspection_id and created_at) for 5 latest inspections from the db
// if the data is fetched successfully calls setInspectionValue method- to set the values for the data objects
function getLatestInspections(user_id) {
	const query = {
	text: 'SELECT id, folder_id, created_at, template_name FROM inspections WHERE user_id = $1 AND folder_id IN ( SELECT portfolio_id FROM portfolio_access_controls WHERE user_id = $2) AND archived_at IS NULL ORDER BY created_at DESC LIMIT 5',
	values: [user_id,user_id],
	};
	client.query(query, (err, res) => {
	if(err){
		throw err;
	}
	else{
		setInspectionValue(res);
	}
	})
}
// sets the values for the inspection_data object
function setInspectionValue(value) {
  for (var i = 0; i < value.rowCount; i++) {
        inspection_data.i_id.push(value.rows[i]['id']);
        inspection_data.i_folder_id.push(value.rows[i]['folder_id']);
        inspection_data.i_created_at.push(value.rows[i]['created_at']);
  }
}

// gets the data (public_id and created_at) for 5 latest reports from the db
// if the data is fetched successfully calls setReportValue method- to set the values for the report_data object
function getLatestReports(user_id) {  
  const query = {
  text:'SELECT public_id as id, created_at FROM reports_v3 WHERE user_id = $1 AND folder_id IN ( SELECT portfolio_id FROM portfolio_access_controls WHERE user_id = $2) AND archived_at IS NULL ORDER BY created_at DESC LIMIT 5',
    values: [user_id,user_id],
  };
  client.query(query, (err, res) => {
  if(err){
    throw err;
  }
  else{
    setReportValue(res);
  }
  })
}

//sets the values for the report_data object
function setReportValue(value) {
  for (var i = 0; i < value.rowCount; i++) {
      report_data.r_id.push(value.rows[i]['id']);
      report_data.r_created_at.push(value.rows[i]['created_at']);
  }
}

// checks if the busniess is on an active iap or not
// fetches the expires_at from the db
// calls setIAPStatus
function getIAPStatus(user_id) {
  const query = {
  text:'SELECT expires_at FROM iap_receipts WHERE company_id IN (SELECT business_id FROM business_membership WHERE user_id = $1)',
    values: [user_id],
  };

  client.query(query, (err, res) => {
  if(err){
    throw err;
  }
  else{
    setIAPStatus(res);
  }
  })
}

// Setting the value for the iap_data object
function setIAPStatus(value) {
  if (value.rowCount!=0) {
        iap_data.iap_expired_at.push(value.rows[0]['expires_at']);
    }
}

// fetching any business notes from the admin added by HappyCo
// calls set method
function getBusinessNotes(user_id) {
  const query = {
  text:'SELECT body FROM admin_notes WHERE noteable_id IN (SELECT business_id FROM business_membership WHERE user_id = $1)',
    values: [user_id],
  };

  client.query(query, (err, res) => {
  if(err){
    throw err;
  }
  else{
    SetBusinessNotes(res);
  }
      client.end()
  })
}

// sets the business notes to the object
function SetBusinessNotes(value) {
  if (value.rowCount!=0) {
    business_notes_data.b_body.push(value.rows[0]['body']);
    }
}
// builds the note to display in the intercom conversation
function noteBuilder(){
  message = '<b>A small note from Yumi 🐶</b><br/><br/>'
  message += '<b>✅   Yumi found these recent <em>Inspections:</em></b><br/>'
  
  //inspection message
  var inspection_id = inspection_data['i_id'];
  var inspection_folder_id = inspection_data['i_folder_id'];
  var inspection_created_at = inspection_data['i_created_at'];

  for(var i in inspection_id){
       var url = "https://manage.happyco.com/folder/"+inspection_folder_id[i]+"/inspections/"+inspection_id[i];
       var date = inspection_created_at[i];
      message += "<a href="+url+">"+url+"</a>";
  }
     message += "<br/><br/>";
  
  //report message
  var report_id = report_data['r_id'];
  var report_created_at = report_data['r_created_at'];

  for(var i in report_id){
       var url = "https://manage.happyco.com/reports/"+report_id[i];
       var date = report_created_at[i];
      message += "<a href="+url+"target=\"_blank\">"+url+"</a>";
  }
     message += "<br/><br/>";
  //optional iap message
  if(iap_data.iap_expired_at.length==1)
    message += "IAP expires on "+iap_data.iap_expired_at;  console.log(message);

  message += "<br/><br/>";

  //optional business notes message  
  if (business_notes_data.b_body.length>=1) {
    message +="Business notes from admin: "+ business_notes_data.b_body;
  }
}
app.listen(3000, () => console.log('Happy app listening on port 3000'));
//module.exports = app;
