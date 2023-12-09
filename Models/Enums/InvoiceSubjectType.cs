using System;
namespace SendTaxDataApp.Models.Enums{
public enum InvoiceSubjectType{
     [System.ComponentModel.Description("اصلی")]
    Main=1,
     [System.ComponentModel.Description("اصلاحی")]
    Modified=2,
     [System.ComponentModel.Description("ابطالی")]
    Canceled=3,
     [System.ComponentModel.Description("برگشت از فروش")]
    ReturnFromSale=4

}}