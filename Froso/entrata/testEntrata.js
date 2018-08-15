  const puppeteer = require('puppeteer');

(async () => {
  
  const browser = await puppeteer.launch({
        ignoreHTTPSErrors: true
    });
  const page = await browser.newPage();

if(typeof require !== 'undefined') XLSX = require('xlsx');

  var workbook = XLSX.readFile('Entrata.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
  
  var address_of_cell = 'A2';
  var desired_cell = worksheet[address_of_cell];
  var desired_value = (desired_cell ? desired_cell.v : undefined);

  for (var i = 13; i < 1339; i++) {
    console.log("Row: "+i);
    
    address_of_cell = 'A'+i
    desired_cell = worksheet[address_of_cell];
    desired_value = (desired_cell ? desired_cell.v : undefined);
    var name = "";
    var url = "https://"+desired_value
    
    var marketCellAddress = 'C'+i
    var marketCell = worksheet[marketCellAddress];
    var marketValue = (marketCell ? marketCell.v : undefined);

    await page.goto(url); 
    
    if (await page.evaluate(() => document.querySelector('.widget.corporate-logo')) != null) {
      name = await page.evaluate(() => document.querySelector('.widget.corporate-logo > a').href)
      console.log(name); 
      if (name.includes("winncompanies.com")) { 
        console.log("Winn Residential");
       
        var writeCell = 'B'+i
        worksheet[writeCell].v = "Winn Residential";
        
        var writeCellMarket = 'C'+i
        worksheet[writeCellMarket].v = "Enterprise";

        continue;
      }
    }
    XLSX.writeFile(workbook ,'Entrata2.xlsx')
  }
  await browser.close();
})();