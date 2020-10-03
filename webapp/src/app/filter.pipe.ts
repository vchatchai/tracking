import { Pipe, PipeTransform } from '@angular/core';
import { formatDate } from '@angular/common';
import { State } from './book.service';


@Pipe({ name: 'appFilter' })
export class FilterPipe implements PipeTransform {
    /**
     * Transform
     *
     * @param {any[]} items
     * @param {string} searchText
     * @returns {any[]}
     */
    transform(items: any[], searchText: string, rowData: any[]): any[] {

        if (!items) {
            return [];
        }
        if (!searchText) {
            return items;
        }
        searchText = searchText.toLowerCase(); 

        let result = false;
        for(let it of items) {
            if (!result) {
                if (it.field == 'no') {
                    rowData['no'] = 'no'
                }else if (it.field != 'size' && it.field != 'gate_out_date'&& it.field != 'gate_in_date') {
                    if (rowData[it.field] != null) {
                        if (!result) {
                            result = rowData[it.field].toLocaleLowerCase().includes(searchText);
                        }
                    }

                } else if (it.field == 'size') {
                    if (rowData[it.field] != null) {
                        let value = `${rowData[it.field]}/${rowData['type']}`
                        if (!result) {
                            result = value.toLocaleLowerCase().includes(searchText)
                            // result = rowData[it.field] == searchText || rowData['type'].toLocaleLowerCase().includes(searchText);
                        }
                    }
                } else if (it.field == 'gate_out_date') {
                    if (!result) {
                        result = formatDate(rowData[it.field], 'dd/MM/yyyy', 'en-US').toLocaleLowerCase().includes(searchText);
                    }
                } else if (it.field == 'gate_in_date') {
                    if (!result) {
                        result = formatDate(rowData[it.field], 'dd/MM/yyyy', 'en-US').toLocaleLowerCase().includes(searchText);
                    }
                }
            }else {
                break;
            }
        }
        return items.filter(it => {
            return result;
        });

    }
}