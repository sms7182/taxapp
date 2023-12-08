using System;
namespace SendTaxDataApp.Models{

public class InvoiceItem{
    private decimal fee;
    private decimal am;
    private decimal dis;
    private decimal vra;
    public string Sstid { get; set; }
    public decimal Am{ get{return am;} set{am=value; Calculate();} }
    public decimal Fee { get{return fee;} set{fee=value;Calculate();} }
    public  decimal Prdis{get;set;}
    public decimal Dis{get{return dis;}set{dis=value;Calculate();}}
    public decimal Adis { get; set; }
    public decimal Vra { get{return vra;} set{vra=value;Calculate();} }
    public decimal Vam{get;set;}
    public decimal Tsstam{get;set;}
    private void Calculate(){
        Prdis=fee*am;
        Adis=Prdis-dis;
        Vam=(vra*Adis/100);
        Tsstam=Adis+Vam;
    }
}}