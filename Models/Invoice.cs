using System;
namespace SendTaxDataApp.Models
{
    

public class Invoice{
public int Setm { get; set; }
public int Inp { get; set; }
public int Ins { get; set; }
public int Inty { get; set; }
public string Tins { get; set; }    
public string Tinb { get; set; }
public decimal Tadis { get; set; }
public decimal Tbill { get; set; }
public decimal Tvam { get; set; }

public decimal Tprdis { get; set; }
public decimal Todam { get; set; }
public decimal Tdis {get;set;}
public decimal Cap { get; set; }  
private DateTime invoiceDate;
public double Indatim{get;set;}
public string Username { get; set; }="A16XAX";

public List<InvoiceItem> Detail{get;set;}
}
}