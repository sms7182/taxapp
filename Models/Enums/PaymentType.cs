using System;
namespace SendTaxDataApp.Models.Enums{
public enum PaymentType{
     [System.ComponentModel.Description("نقدی")]
    Cashe=1,
     [System.ComponentModel.Description("نسیه")]
    Credit=2,
     [System.ComponentModel.Description("نقدی/نسیه")]
    CasheCredit=3

}}