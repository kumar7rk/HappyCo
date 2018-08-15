  const puppeteer = require('puppeteer');

//  const {performance} = require('perf_hooks');

(async () => {
  
  const browser = await puppeteer.launch({
          // headless: false
          // ignoreHTTPSErrors: true
    });
  const page = await browser.newPage();
  //var sound = new Audio('/audio.mp3');
  if(typeof require !== 'undefined') XLSX = require('xlsx');

  var workbook = XLSX.readFile('Entrata.xlsx');
  var first_sheet_name = workbook.SheetNames[0];
  var worksheet = workbook.Sheets[first_sheet_name];
  
  var address_of_cell = 'A2';
  var desired_cell = worksheet[address_of_cell];
  var desired_value = (desired_cell ? desired_cell.v : undefined);

  // var t0 = performance.now();

  for (var i = 1338; i < 1339; i++) {
    console.log("Row: "+i);
    
    address_of_cell = 'A'+i
    desired_cell = worksheet[address_of_cell];
    desired_value = (desired_cell ? desired_cell.v : undefined);
    var name = "";
    var url = "https://"+desired_value
    
    var marketCellAddress = 'C'+i
    var marketCell = worksheet[marketCellAddress];
    var marketValue = (marketCell ? marketCell.v : undefined);
    try{
      await page.goto(url); 
    }
    catch(error){
      console.log(error);
      continue;
    }
  try{
    if (await page.evaluate(() => document.querySelector('.corporate_logo.corporate-logo-item > a')) != null) {
      name = await page.evaluate(() => document.querySelector('.corporate_logo.corporate-logo-item > a').href)
      
      if (name.includes("trinity-pm.com")) { //16
        console.log("Trinity");
        
        var writeCell = 'B'+i
        worksheet[writeCell].v = "Trinity Property Consultants";
        
        var writeCellMarket = 'C'+i
        worksheet[writeCellMarket].v = "Enterprise";

        continue;
      }
      else{
            var writeCell = 'B'+i
            worksheet[writeCell].v = name;
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Website";
            continue;
      }
    }
    
    if (await page.evaluate(() => document.querySelector('.footer-copyright')) != null) {
      name = await page.evaluate(() => document.querySelector('.footer-copyright').textContent)
      
      if (name.includes("Greystar")) { //104
        console.log("Greystar");
       
        var writeCell = 'B'+i
        worksheet[writeCell].v = "Greystar";
        
        var writeCellMarket = 'C'+i
        worksheet[writeCellMarket].v = "Enterprise";

        continue;
      }
    }
    
    if (await page.evaluate(() => document.querySelector('.corporate_logo_1.corporate-logo-item > a')) != null) {
          name = await page.evaluate(() => document.querySelector('.corporate_logo_1.corporate-logo-item > a').href)
          name = name.toLowerCase();
          if (name.includes("fairfieldresidential.com")) { //105
            console.log("Fairfield");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Fairfield Residential";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";

            continue;
          }
          else if (name.includes("ca-studentliving.com")) { //15
            console.log("CA Student Living");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "CA Student Living";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("hanoverco.com")) { //487
            console.log("Hanover Company");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Hanover Company";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("lincolnapts.com")) { //122
            console.log("Lincoln Property Company");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Lincoln Property Company";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("kettler.com")) { //
            console.log("Kettler");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Kettler";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("fogelman-management.com")) { //
            console.log("Fogelman Management Group");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Fogelman Management Group";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("csmcorp.net")) { //
            console.log("CSM Corp");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "CSM Corp";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("tlcproperties.com")) { //
            console.log("TLC Properties");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "TLC Properties";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("caliberliving.com")) { //
            console.log("Caliber Living");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Caliber Living";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("milestonerents.com")) { //
            console.log("Milestone Management");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Milestone Management";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("fpimgt.com")) { //
            console.log("FPI Management");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "FPI Management";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("jrkpropholdings.com")) { //
            console.log("JRK");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "JRK";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("fogelman-management.com")) { //
            console.log("Fogelman Management");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Fogelman Management";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("baronproperties.com")) { //
            console.log("Baron Properties");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Baron Properties";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("cambridgemgmt.net")) { //
            console.log("Cambridge Management");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Cambridge Management";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("echelonpg.com")) { //
            console.log("Echelon Property Group");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Echelon Property Group";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("woodruffre.com")) { //
            console.log("Woodruff");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Woodruff";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          }
          else if (name.includes("stevebrownapts.com")) { //
            console.log("Steve Brown Apts");
            
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Steve Brown Apts";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";
            
            continue;
          } 
          else if (name.includes("millcreekplaces.com")) { //54
            console.log("Mill creek");
          
            var writeCell = 'B'+i
            worksheet[writeCell].v = "Mill Creek";
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Enterprise";

            continue;
          }
          else{
            console.log(name +"Fairfield")
            var writeCell = 'B'+i
            worksheet[writeCell].v = name;
            
            var writeCellMarket = 'C'+i
            worksheet[writeCellMarket].v = "Website";

            continue;
          }
      } 
    
    if (await page.evaluate(() => document.querySelector('.corporate-logo-container.corporate-logo > a')) != null) {
      name = await page.evaluate(() => document.querySelector('.corporate-logo-container.corporate-logo > a').href)
      
      if (name.includes("millcreekplaces.com")) { //54
        console.log("Mill creek");
      
        var writeCell = 'B'+i
        worksheet[writeCell].v = "Mill Creek";
        
        var writeCellMarket = 'C'+i
        worksheet[writeCellMarket].v = "Enterprise";

        continue;
      }
    }

    if (await page.evaluate(() => document.querySelector('.widget.corporate-logo')) != null) {
      name = await page.evaluate(() => document.querySelector('.widget.corporate-logo > a').href)
      
      if (name.includes("winncompanies.com")) { 
        console.log("Winn Residential");
       
        var writeCell = 'B'+i
        worksheet[writeCell].v = "Winn Residential";
        
        var writeCellMarket = 'C'+i
        worksheet[writeCellMarket].v = "Enterprise";

        continue;
      }
    }
    XLSX.writeFile(workbook ,'Entrata5.xlsx')
  }
  catch(error){
    console.log(error);
    continue;
  }
  }
  // var t1 = performance.now();
  //console.log((t1-t0)/1000+ " seconds");
  
  //sound.play();

    /*var snd2 = new Audio("data:audio/wav;base64,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qJrx4KOQjahZxWbcZgztj2c49nKmkId44S71j0c8eV9yDK6uPRzx5X18eDvjvQ6yKo9ZSS6l//8elePK/Lf//IInrOF/FvDoADYAGBMGb7FtErm5MXMlmPAJQVgWta7Zx2go+8xJ0UiCb8LHHdftWyLJE0QIAIsI+UbXu67dZMjmgDGCGl1H+vpF4NSDckSIkk7Vd+sxEhBQMRU8j/12UIRhzSaUdQ+rQU5kGeFxm+hb1oh6pWWmv3uvmReDl0UnvtapVaIzo1jZbf/pD6ElLqSX+rUmOQNpJFa/r+sa4e/pBlAABoAAAAA3CUgShLdGIxsY7AUABPRrgCABdDuQ5GC7DqPQCgbbJUAoRSUj+NIEig0YfyWUho1VBBBA//uQZB4ABZx5zfMakeAAAAmwAAAAF5F3P0w9GtAAACfAAAAAwLhMDmAYWMgVEG1U0FIGCBgXBXAtfMH10000EEEEEECUBYln03TTTdNBDZopopYvrTTdNa325mImNg3TTPV9q3pmY0xoO6bv3r00y+IDGid/9aaaZTGMuj9mpu9Mpio1dXrr5HERTZSmqU36A3CumzN/9Robv/Xx4v9ijkSRSNLQhAWumap82WRSBUqXStV/YcS+XVLnSS+WLDroqArFkMEsAS+eWmrUzrO0oEmE40RlMZ5+ODIkAyKAGUwZ3mVKmcamcJnMW26MRPgUw6j+LkhyHGVGYjSUUKNpuJUQoOIAyDvEyG8S5yfK6dhZc0Tx1KI/gviKL6qvvFs1+bWtaz58uUNnryq6kt5RzOCkPWlVqVX2a/EEBUdU1KrXLf40GoiiFXK///qpoiDXrOgqDR38JB0bw7SoL+ZB9o1RCkQjQ2CBYZKd/+VJxZRRZlqSkKiws0WFxUyCwsKiMy7hUVFhIaCrNQsKkTIsLivwKKigsj8XYlwt/WKi2N4d//uQRCSAAjURNIHpMZBGYiaQPSYyAAABLAAAAAAAACWAAAAApUF/Mg+0aohSIRobBAsMlO//Kk4soosy1JSFRYWaLC4qZBYWFRGZdwqKiwkNBVmoWFSJkWFxX4FFRQWR+LsS4W/rFRb/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////VEFHAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAU291bmRib3kuZGUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMjAwNGh0dHA6Ly93d3cuc291bmRib3kuZGUAAAAAAAAAACU=");  
  
    snd2.play();*/

  await browser.close();
})();